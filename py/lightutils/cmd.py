import os
import sys
import json
import argparse
import ruamel.yaml
from .client import Client


class Cli(object):

    def __init__(self, parser):
        self._parser = parser
        self._subparsers = {}

    def parser(self, name, **kws):
        def decorate(func):
            sub = self._subparsers.get(name)
            if sub is None:
                sub = self._parser.add_parser(name, **kws)
                self._subparsers[name] = sub
            sub.set_defaults(func=func)
            return func

        return decorate

    def argument(self, name, *args,  **kws):
        def decorate(func):
            sub = self._subparsers.get(name)
            if sub is not None:
                sub.add_argument(*args, **kws)
            return func

        return decorate

    def output(self, func):
        def decorate(opt, *args, **kws):
            resp = func(opt, *args, **kws)
            if opt.output == 'json':
                json.dump(resp.json(), sys.stdout, indent=2)
            elif opt.output == 'yaml':
                yaml = ruamel.yaml.YAML()
                yaml.dump(resp.json(), sys.stdout)
            else:
                sys.stdout.write(resp.text)
            return resp

        return decorate


def with_client(func):
    def decorate(opt, *args, **kws):
        c = Client(opt.ls_url, opt.ls_auth, debug=opt.debug)
        return func(c, opt, *args, **kws)

    return decorate


parse = argparse.ArgumentParser()
parse.add_argument('-u', '--ls-auth',
                   help='set basic auth',
                   default=os.environ.get("LS_AUTH", ""))
parse.add_argument('-l', '--ls-url',
                   help='set server url',
                   default=os.environ.get("LS_URL", "https://localhost:10080"))
parse.add_argument('-d', '--debug',
                   help='print verbose log',
                   action='store_true', default=False)
parse.add_argument('-o', '--output',
                   help='set output format such as json, yaml',
                   default='yaml')


def parse_args():
    args = parse.parse_args()
    return args, parse


cli = Cli(parse.add_subparsers())


@cli.argument('list', "-p", "--path", help="list data of path")
@cli.parser('list', help="display data")
@cli.output
@with_client
def cmd_list_index(client, opt):
    if opt.path is None:
        opt.path = '/api/hyper'
    return client.request(opt.path, "GET")

