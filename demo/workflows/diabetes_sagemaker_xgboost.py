from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

import joblib
import pandas as pd
from flytekit.sdk.tasks import python_task, outputs, inputs
from flytekit.sdk.types import Types
from flytekit.sdk.workflow import workflow_class, Output, Input
from sklearn.metrics import accuracy_score
from sklearn.model_selection import train_test_split
from xgboost import XGBClassifier

import workflows.diabetes_xgboost as dxgb

@workflow_class
class DiabetesXGBoostModelOptimizer(object):
    """
    This pipeline trains an XGBoost mode for any given dataset that matches the schema as specified in
    https://github.com/jbrownlee/Datasets/blob/master/pima-indians-diabetes.names.
    """

    # Inputs dataset, fraction of the dataset to be split out for validations and seed to use to perform the split
    dataset = Input(Types.CSV, default=Types.CSV.create_at_known_location(
        "https://raw.githubusercontent.com/jbrownlee/Datasets/master/pima-indians-diabetes.data.csv"),
                    help="A CSV File that matches the format https://github.com/jbrownlee/Datasets/blob/master/pima-indians-diabetes.names")

    test_split_ratio = Input(Types.Float, default=0.33, help="Ratio of how much should be test to Train")
    seed = Input(Types.Integer, default=7, help="Seed to use for splitting.")

    # the actual algorithm
    split = dxgb.get_traintest_splitdatabase(dataset=dataset, seed=seed, test_split_ratio=test_split_ratio)
    fit_task = dxgb.fit(x=split.outputs.x_train, y=split.outputs.y_train, hyperparams=dxgb.XGBoostModelHyperparams(max_depth=4).to_dict())
    predicted = dxgb.predict(model_ser=fit_task.outputs.model, x=split.outputs.x_test)
    score_task = dxgb.metrics(predictions=predicted.outputs.predictions, y=split.outputs.y_test)

    # Outputs: joblib seralized model and accuracy of the model
    model = Output(fit_task.outputs.model, sdk_type=Types.Blob)
    accuracy = Output(score_task.outputs.accuracy, sdk_type=Types.Float)
