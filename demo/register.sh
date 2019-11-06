PROJECT=aws
DOMAIN=development

pyflyte -p ${PROJECT} -d ${DOMAIN} --config flyte.config register workflows
