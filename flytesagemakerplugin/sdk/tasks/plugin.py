from __future__ import absolute_import

import datetime as _datetime
from flytekit.sdk import types as _sdk_types
from flytekit.common.tasks import task as _sdk_task
from flytekit.common import interface as _interface
from flytekit.models import literals as _literal_models, task as _task_models, interface as _interface_model


_TASK_TYPE = "aws-sagemaker"


class SagemakerXgBoostTrainer(_sdk_task.SdkTask):

    def __init__(
        self,
        role_arn,
        resource_config,
        algorithm_specification=None,
        stopping_condition=None,
        retries=0,
        cacheable=False,
        cache_version="",
    ):
        """
        :param Text role_arn: This probably shouldn't be statically configured here.
        :param dict[Text,T] algorithm_specification: https://docs.aws.amazon.com/sagemaker/latest/dg/API_AlgorithmSpecification.html
        :param dict[Text,T] resource_config: https://docs.aws.amazon.com/sagemaker/latest/dg/API_ResourceConfig.html
        :param dict[Text,T] stopping_condition: https://docs.aws.amazon.com/sagemaker/latest/dg/API_StoppingCondition.html
        :param int retries: Number of time to retry.
        :param bool cacheable: Whether or not to use Flyte's caching system.
        :param Text cache_version: Update this to notify a behavioral change requiring the cache to be invalidated.
        """

        # TODO: We probably don't want to hardcode things here, but for example...
        algorithm_specification = algorithm_specification or {}
        algorithm_specification["TrainingImage"] = algorithm_specification.get("TrainingImage") or \
            "825641698319.dkr.ecr.us-east-2.amazonaws.com/xgboost:1"
        algorithm_specification["TrainingInputMode"] = "File"
        resource_config = resource_config or {}
        stopping_condition = stopping_condition or {}

        # TODO: Optionally, pull timeout behavior from stopping condition and pass to Flyte task def.
        timeout = _datetime.timedelta(seconds=0)

        # TODO: The FlyteKit type engine is extensible so we can create a SagemakerInput type with custom
        # TODO:     parsing/casting logic. For now, we will use the Generic type since there is a little that needs
        # TODO:     to be done on Flyte side to unlock this cleanly.
        # TODO: This call to the super-constructor will be less verbose in future versions of Flytekit following a
        # TODO:     refactor.
        # TODO: Add more configurations to the custom dict. These are things that are necessary to execute the task,
        # TODO:     but might not affect the outputs (i.e. Running on a bigger machine). These are currently static for
        # TODO:     a given definition of a task, but will be more dynamic in the future. Also, it is possible to
        # TODO:     make it dynamic by using our @dynamic_task.
        # TODO: You might want to inherit the role ARN from the execution at runtime.
        super(SagemakerXgBoostTrainer, self).__init__(
            type=_TASK_TYPE,
            metadata=_task_models.TaskMetadata(
                discoverable=cacheable,
                runtime=_task_models.RuntimeMetadata(0, "0.1.0b0", "sagemaker"),
                timeout=timeout,
                retries=_literal_models.RetryStrategy(retries=retries),
                discovery_version=cache_version,
                deprecated_error_message="",
            ),
            interface=_interface.TypedInterface({}, {}),
            custom={
                "RoleArn": role_arn,
                "AlgorithmSpecification": algorithm_specification,
                "ResourceConfig": resource_config,
                "StoppingCondition": stopping_condition,
            }
        )

        # TODO: Add more inputs that we expect to change the outputs of the task.
        # TODO: We can add outputs too!
        # We use helper methods for adding to interface, thus overriding the one set above. This will be simplified post
        # refactor.
        self.add_inputs(
            {
                'static_hyperparameters':
                    _interface_model.Variable(_sdk_types.Types.Generic.to_flyte_literal_type(), ''),
                'train': _interface_model.Variable(_sdk_types.Types.MultiPartCSV.to_flyte_literal_type(), ''),
                'validation': _interface_model.Variable(_sdk_types.Types.MultiPartCSV.to_flyte_literal_type(), ''),
            }
        )
        self.add_outputs(
            {
                'model': _interface_model.Variable(_sdk_types.Types.Blob.to_flyte_literal_type(), '')
            }
        )
