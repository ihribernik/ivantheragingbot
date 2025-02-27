from setuptools import find_packages, setup

with open("README.md", "r", encoding="utf-8") as file:
    long_description = file.read()

setup(
    name="ivantheraginbot",
    version="0.1",
    packages=find_packages(),
    install_requires=["requests", "python-dotenv"],
    long_description=long_description,
)
