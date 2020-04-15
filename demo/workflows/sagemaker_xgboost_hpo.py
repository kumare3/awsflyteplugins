from __future__ import absolute_import
from __future__ import division
from __future__ import print_function

import tarfile
import pandas as pd
import pickle
from flytekit.common import utils
from flytekit.sdk.tasks import python_task, outputs, inputs
from flytekit.sdk.types import Types
from flytekit.sdk.workflow import workflow_class, Input, Output
from flytesagemakerplugin.sdk.tasks.plugin import SagemakerXgBoostOptimizer
from xgboost import XGBClassifier

example_hyperparams = {
    "base_score": "0.5",
    "booster": "gbtree",
    "csv_weights": "0",
    "dsplit": "row",
    "grow_policy": "depthwise",
    "lambda_bias": "0.0",
    "max_bin": "256",
    "max_leaves": "0",
    "normalize_type": "tree",
    "objective": "reg:linear",
    "one_drop": "0",
    "prob_buffer_row": "1.0",
    "process_type": "default",
    "rate_drop": "0.0",
    "refresh_leaf": "1",
    "sample_type": "uniform",
    "scale_pos_weight": "1.0",
    "silent": "0",
    "sketch_eps": "0.03",
    "skip_drop": "0.0",
    "tree_method": "auto",
    "tweedie_variance_power": "1.5",
    "updater": "grow_colmaker,prune",
}

xgtrainer_task = SagemakerXgBoostOptimizer(
    region="us-east-1",
    role_arn="arn:aws:iam::173840052742:role/modelbuilderapibatchworker-staging",
    resource_config={
        "InstanceCount": 1,
        "InstanceType": "ml.m4.xlarge",
        "VolumeSizeInGB": 25,
    },
    stopping_condition={"MaxRuntimeInSeconds": 43200, "MaxWaitTimeInSeconds": 43200},
    algorithm_specification={"TrainingImage": "811284229777.dkr.ecr.us-east-1.amazonaws.com/xgboost:latest",
                             "TrainingInputMode": "File", "AlgorithmName": "xgboost"},
    retries=2,
    cacheable=True,
    cache_version="2.0",
)


def read_and_merge(first, second):
    """
    Sagemaker likes the target to be in column 1. This method takes the y and the x and just places the dataframes
    next to each other, yielding a common dataframe
    """
    with first as r:
        first_df = r.read()
    with second as r:
        second_df = r.read()
    if len(first_df) != len(second_df):
        raise Exception("trying to merge to data frames which are not equal in length")
    return pd.concat([first_df, second_df], axis=1)


@inputs(x_train=Types.Schema(), x_test=Types.Schema(), y_train=Types.Schema(), y_test=Types.Schema())
@outputs(train=Types.MultiPartCSV, validation=Types.MultiPartCSV)
@python_task(cache_version='3.0', cache=True, memory_limit="500Mi")
def convert_to_sagemaker_csv(ctx, x_train, y_train, x_test, y_test, train, validation):
    _train = read_and_merge(y_train, x_train)
    _validate = read_and_merge(y_test, x_test)

    with utils.AutoDeletingTempDir("train") as t:
        f = t.get_named_tempfile("train.csv")
        _train.to_csv(f, header=False, index=False)
        train.set(t.name)

    with utils.AutoDeletingTempDir("validate") as t:
        f = t.get_named_tempfile("validate.csv")
        _validate.to_csv(f, header=False, index=False)
        validation.set(t.name)


@inputs(model_tar=Types.Blob)
@outputs(model=Types.Blob)
@python_task(cache_version="3.0", cache=True, memory_limit="500Mi")
def untar_xgboost(ctx, model_tar, model):
    model_tar.download()
    fname = "xgboost-model"
    with tarfile.open(model_tar.local_path, "r:gz") as tf:
        tf.extract(fname)
    model.set(fname)


@workflow_class
class StructuredSagemakerXGBoostHPO(object):
    # Input parameters
    static_hyperparameters = Input(
        Types.Generic,
        help="A list of the static hyperparameters to pass to the training jobs.",
        default=example_hyperparams,
    )
    train_data = Input(
        Types.Schema(), help="A Columnar schema that contains all the features used for training.",
    )
    train_target = Input(
        Types.Schema(), help="A Columnar schema that contains all the labeled results for train_data.",
    )

    validation_data = Input(
        Types.Schema(), help="A Columnar schema that contains all the features used for validation.",
    )
    validation_target = Input(
        Types.Schema(), help="A Columnar schema that contains all the labeled results for validation_data.",
    )

    sagemaker_transform = convert_to_sagemaker_csv(x_train=train_data, y_train=train_target,
                                                   x_test=validation_data, y_test=validation_target)

    # Node definitions
    train_node = xgtrainer_task(
        static_hyperparameters=static_hyperparameters,
        train=sagemaker_transform.outputs.train,
        validation=sagemaker_transform.outputs.validation,
    )

    untar = untar_xgboost(
        model_tar=train_node.outputs.model,
    )

    # Outputs
    model = Output(untar.outputs.model, sdk_type=Types.Blob)

# Create a launch plan that can be used in other workflows, with default inputs
fit_lp = StructuredSagemakerXGBoostHPO.create_launch_plan()
