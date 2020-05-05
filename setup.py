from __future__ import absolute_import

from setuptools import setup, find_packages  # noqa
import flytesagemakerplugin

setup(
    name='flytesagemakerplugin',
    version=flytesagemakerplugin.__version__,
    maintainer='AWS + Flyte team',
    maintainer_email='ketan.umare@gmail.com',
    packages=find_packages(exclude=["tests*"]),
    url='https://github.com/kumare3/awsflyteplugins',
    description='FlyteKit Example plugin',
    long_description=open('README.md').read(),
    entry_points={
        'console_scripts': [
        ]
    },
    setup_requires=[
        "pip>10.0.0,<19.2.0"
    ],
    install_requires=[
        "protobuf>=3.7.0, <4, !=3.11.*"
        "flytekit>=0.6.0",
    ],
    tests_require=['pytest'],
    scripts=[
    ],
    license="apache2",
    python_requires=">=2.7"
)
