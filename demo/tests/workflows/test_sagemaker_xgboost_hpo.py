from flytekit.sdk.test_utils import flyte_test
from flytekit.sdk.types import Types
from flytekit.common.types.impl import blobs

from workflows import sagemaker_xgboost_hpo as sxh
from workflows import diabetes_xgboost as dxgb

@flyte_test
def test_untar_xgboost():

    b = blobs.Blob.from_python_std("/Users/kumare/go/src/github.com/kumare3/awsflyteplugins/model.tar.gz")
    v = sxh.untar_xgboost.unit_test(model_tar=b)
    print(v)

    dataset = Types.CSV.create_at_known_location(
        "https://raw.githubusercontent.com/jbrownlee/Datasets/master/pima-indians-diabetes.data.csv")

    # Test get dataset
    result = dxgb.get_traintest_splitdatabase.unit_test(dataset=dataset, seed=7, test_split_ratio=0.33)
    assert "x_train" in result
    assert "y_train" in result
    assert "x_test" in result
    assert "y_test" in result

    print(v["model"])
    p = dxgb.predict.unit_test(x=result["x_test"], model_ser=v["model"])