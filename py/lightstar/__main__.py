from .cmd import parse_args


def main():
    args, parse = parse_args()
    if hasattr(args, 'func'):
        args.func(args)
    else:
        parse.print_help()


if __name__ == '__main__':
    main()

