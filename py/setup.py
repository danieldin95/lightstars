from setuptools import setup

setup(
    name='lightstar',
    version='0.0.1',
    author='Daniel Ding',
    author_email='danieldin95@163.com',
    packages=['lightstar'],
    entry_points={
        'console_scripts': [
            'lightstar-utils = lightstar.__main__:main',
        ]
    },
    install_requires=open('requirements.txt').readlines(),
)
