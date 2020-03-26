from setuptools import setup

setup(
    name='lightutils',
    version='0.0.1',
    author='Daniel Ding',
    author_email='danieldin95@163.com',
    packages=['lightutils'],
    entry_points={
        'console_scripts': [
            'lightutils = lightutils.__main__:main',
        ]
    },
    install_requires=open('requirements.txt').readlines(),
)
