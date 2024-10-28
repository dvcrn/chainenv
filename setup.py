
from setuptools import setup, find_packages

setup(
    name="chainenv",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[],
    entry_points={
        'console_scripts': [
            'chainenv=cli.cli:main',
        ],
    },
    author="David Mohl",
    author_email="git@d.sh",
    description="Wrapper around macOS Keychain to quickly set/get passwords",
    long_description=open("README.md").read(),
    long_description_content_type="text/markdown",
    url="https://github.com/dvcrn/chainenv",
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: MacOS :: MacOS X",
    ],
    python_requires=">=3.6",
)
