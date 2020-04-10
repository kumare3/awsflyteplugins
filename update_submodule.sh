DIR=$(dirname "${BASH_SOURCE[0]}")
pushd $DIR
git submodule update --init
pushd amazon-sagemaker-operator-for-k8s
git fetch
git checkout origin/r1.1
CURRENT_TAG=$(git describe --abbrev=0 --tags)
popd
popd
echo "Version Updated to ${CURRENT_TAG}, use this tag when you release awsflyteplugins"
