.. _api_file_sagemaker.proto:

sagemaker.proto
===============

.. _api_msg_flyte.plugins.sagemaker.AlgorithmSpecification:

flyte.plugins.sagemaker.AlgorithmSpecification
----------------------------------------------

`[flyte.plugins.sagemaker.AlgorithmSpecification proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L6>`_


.. code-block:: json

  {
    "TrainingImage": "...",
    "TrainingInputMode": "...",
    "AlgorithmName": "...",
    "MetricDefinitions": []
  }

.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.TrainingImage:

TrainingImage
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.TrainingInputMode:

TrainingInputMode
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.AlgorithmName:

AlgorithmName
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinitions:

MetricDefinitions
  (:ref:`flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition <api_msg_flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition>`) 
  
.. _api_msg_flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition:

flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition
---------------------------------------------------------------

`[flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L11>`_


.. code-block:: json

  {
    "Name": "...",
    "Regex": "..."
  }

.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition.Name:

Name
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.AlgorithmSpecification.MetricDefinition.Regex:

Regex
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  



.. _api_msg_flyte.plugins.sagemaker.ResourceConfig:

flyte.plugins.sagemaker.ResourceConfig
--------------------------------------

`[flyte.plugins.sagemaker.ResourceConfig proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L18>`_


.. code-block:: json

  {
    "InstanceType": "...",
    "InstanceCount": "...",
    "VolumeSizeInGB": "...",
    "VolumeKmsKeyId": "..."
  }

.. _api_field_flyte.plugins.sagemaker.ResourceConfig.InstanceType:

InstanceType
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.ResourceConfig.InstanceCount:

InstanceCount
  (`int64 <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.ResourceConfig.VolumeSizeInGB:

VolumeSizeInGB
  (`int64 <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.ResourceConfig.VolumeKmsKeyId:

VolumeKmsKeyId
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  


.. _api_msg_flyte.plugins.sagemaker.StoppingCondition:

flyte.plugins.sagemaker.StoppingCondition
-----------------------------------------

`[flyte.plugins.sagemaker.StoppingCondition proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L25>`_


.. code-block:: json

  {
    "MaxRuntimeInSeconds": "...",
    "MaxWaitTimeInSeconds": "..."
  }

.. _api_field_flyte.plugins.sagemaker.StoppingCondition.MaxRuntimeInSeconds:

MaxRuntimeInSeconds
  (`int64 <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.StoppingCondition.MaxWaitTimeInSeconds:

MaxWaitTimeInSeconds
  (`int64 <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  


.. _api_msg_flyte.plugins.sagemaker.VpcConfig:

flyte.plugins.sagemaker.VpcConfig
---------------------------------

`[flyte.plugins.sagemaker.VpcConfig proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L30>`_


.. code-block:: json

  {
    "SecurityGroupIds": [],
    "Subnets": []
  }

.. _api_field_flyte.plugins.sagemaker.VpcConfig.SecurityGroupIds:

SecurityGroupIds
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.VpcConfig.Subnets:

Subnets
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  


.. _api_msg_flyte.plugins.sagemaker.SagemakerHPOJob:

flyte.plugins.sagemaker.SagemakerHPOJob
---------------------------------------

`[flyte.plugins.sagemaker.SagemakerHPOJob proto] <https://github.com/lyft/flyteidl/blob/master/protos/sagemaker.proto#L35>`_


.. code-block:: json

  {
    "RoleArn": "...",
    "AlgorithmSpecification": "{...}",
    "ResourceConfig": "{...}",
    "StoppingCondition": "{...}",
    "VpcConfig": "{...}",
    "EnableSpotTraining": "..."
  }

.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.RoleArn:

RoleArn
  (`string <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  
.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.AlgorithmSpecification:

AlgorithmSpecification
  (:ref:`flyte.plugins.sagemaker.AlgorithmSpecification <api_msg_flyte.plugins.sagemaker.AlgorithmSpecification>`) 
  
.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.ResourceConfig:

ResourceConfig
  (:ref:`flyte.plugins.sagemaker.ResourceConfig <api_msg_flyte.plugins.sagemaker.ResourceConfig>`) 
  
.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.StoppingCondition:

StoppingCondition
  (:ref:`flyte.plugins.sagemaker.StoppingCondition <api_msg_flyte.plugins.sagemaker.StoppingCondition>`) 
  
.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.VpcConfig:

VpcConfig
  (:ref:`flyte.plugins.sagemaker.VpcConfig <api_msg_flyte.plugins.sagemaker.VpcConfig>`) 
  
.. _api_field_flyte.plugins.sagemaker.SagemakerHPOJob.EnableSpotTraining:

EnableSpotTraining
  (`bool <https://developers.google.com/protocol-buffers/docs/proto#scalar>`_) 
  

