import json
import os
import sys
from flytesagemakerplugin.sdk.tasks.plugin import SagemakerXgBoostOptimizer
from flytekit.sdk.workflow import workflow_class, Input, Output
from flytekit.sdk.types import Types
from flytekit.configuration import TemporaryConfiguration

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
    region="us-east-2",
    role_arn="arn:aws:iam::123456789012:role/service-role/AmazonSageMaker-ExecutionRole",
    resource_config={
        "InstanceCount": 1,
        "InstanceType": "ml.m4.xlarge",
        "VolumeSizeInGB": 25,
    },
    stopping_condition={"MaxRuntimeInSeconds": 43200, "MaxWaitTimeInSeconds": 43200},
    algorithm_specification={"TrainingInputMode": "File", "AlgorithmName": "xgboost"},
    retries=2,
)


@workflow_class
class DemoWorkflow(object):
    # Input parameters
    static_hyperparameters = Input(
        Types.Generic,
        help="A list of the static hyperparameters to pass to the training jobs.",
    )
    train_data = Input(
        Types.MultiPartCSV, help="S3 path to a flat directory of CSV files."
    )
    validation_data = Input(
        Types.MultiPartCSV, help="S3 path to a flat directory of CSV files."
    )

    # Node definitions
    train_node = xgtrainer_task(
        static_hyperparameters=example_hyperparams,
        train=train_data,
        validation=validation_data,
    )

    # Outputs
    trained_model = Output(train_node.outputs.model, sdk_type=Types.Blob)


if __name__ == "__main__":
    _PROJECT = "aws"
    _DOMAIN = "development"
    _USAGE = (
        "Usage:\n\n"
        "\tpython sample.py render_task\n"
        "\tpython sample.py execute <version> <train data path> <validation data path> <hyperparameter json>\n"
    )

    with TemporaryConfiguration(
        os.path.join(os.path.dirname(__file__), "..", "flyte.config")
    ):
        if sys.argv[1] == "render_task":
            print("Task Definition:\n\n")
            print(xgtrainer_task.to_flyte_idl())
            print("\n\n")
        elif sys.argv[1] == "execute":
            if len(sys.argv) != 6:
                print(_USAGE)
            else:
                try:
                    # Register, if not already.
                    xgtrainer_task.register(
                        _PROJECT, _DOMAIN, "xgtrainer_task", sys.argv[2]
                    )
                    DemoWorkflow.register(
                        _PROJECT, _DOMAIN, "DemoWorkflow", sys.argv[2]
                    )
                    lp = DemoWorkflow.create_launch_plan()
                    lp.register(_PROJECT, _DOMAIN, "DemoWorkflow", sys.argv[2])
                except:
                    print(
                        "NOTE: If you changed anything about the task or workflow definition, you must register a "
                        "new unique version."
                    )
                    raise
                ex = lp.execute(
                    _PROJECT,
                    _DOMAIN,
                    inputs={
                        "train_data": sys.argv[3],
                        "validation_data": sys.argv[4],
                        "static_hyperparameters": example_hyperparams,
                    },
                )
                print("Waiting for execution to complete...")
                ex.wait_for_completion()
                ex.sync()
                print(
                    "Trained model is available here: {}".format(
                        ex.outputs.trained_model.uri
                    )
                )
        else:
            print(_USAGE)
