import json
import os
import sys
from flytesagemakerplugin.sdk.tasks.plugin import SagemakerXgBoostOptimizer
from flytekit.sdk.workflow import workflow_class, Input, Output
from flytekit.sdk.types import Types
from flytekit.configuration import TemporaryConfiguration


xgtrainer_task = SagemakerXgBoostOptimizer(
    role_arn="TODO: Fill me in",
    resource_config={
        "InstanceCount": 1,
        "InstanceType": "ml.c4.large",
        "VolumeSizeInGB": 10,
    },
    stopping_condition={"MaxRuntimeInSeconds": 43200},
    retries=2,
)


@workflow_class
class DemoWorkflow(object):
    # Input parameters
    train_data = Input(
        Types.MultiPartCSV, help="s3 path to a flat directory of CSV files."
    )
    validation_data = Input(
        Types.MultiPartCSV, help="s3 path to a flat directory of CSV files."
    )

    # Node definitions
    train_node = xgtrainer_task(
        # static_hyperparameters={
        #     "eval_metric": "auc",
        #     "num_round": "100",
        #     "objective": "binary:logistic",
        #     "rate_drop": "0.3",
        #     "tweedie_variance_power": "1.4",
        # },
        train=train_data,
        validation=validation_data,
    )

    # Outputs
    trained_model = Output(train_node.outputs.model, sdk_type=Types.Blob)

    # TODO: Do other things with the resulting model! Ping Matt Smith if you'd like help expanding on this demo
    # for a more holistic example.


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
                    inputs={"train_data": sys.argv[3], "validation_data": sys.argv[4]},
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
