from __future__ import absolute_import

from flytekit.sdk.tasks import outputs
from flytekit.sdk.types import Types
from flytesagemakerplugin.sdk.tasks.plugin import dummy_task

from six.moves import range

# This file is in a subdirectory to make it easier to exclude when not running in a container
# and pyspark is not available


@outputs(out=Types.Integer)
@dummy_task(retries=1)
def my_task(wf, sc, out):
    def _inside(p):
        return p < 1000
    count = sc.parallelize(range(0, 10000)) \
        .filter(_inside).count()
    out.set(count)


@outputs(out=Types.Integer)
@dummy_task(retries=3)
def my_task2(wf, sc, out):
    # This test makes sure spark_task doesn't choke on a non-package module and modules which overlap with auto-included
    # modules.
    def _inside(p):
        return p < 500
    count = sc.parallelize(range(0, 10000)) \
        .filter(_inside).count()
    out.set(count)


def test_basic_task_execution():
    outputs = my_task.unit_test()
    assert outputs['out'] == 1000

    outputs = my_task2.unit_test()
    assert outputs['out'] == 500
