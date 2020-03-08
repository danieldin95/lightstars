export var Utils = {
    // num: int
    iton: function (data, n) {
        return (Array(n).join(0) + data).slice(-n);
    },
    // num: string
    aton: function(data, n) {
        let num = "" + data;
        if (num.length > n) {
            return num
        }
        let ret = "";
        for (let i = 0; i < n - num.length; i++) {
            ret += "0"
        }
        return ret + num;
    },
    toJSON: function (data) {
        let ret = {};
        data.forEach(function (element, index) {
            let key = element['name'];
            let value = element['value'];
            if (key && value) {
                ret[key] = value;
            }
        });
        return ret;
    },
};


export class CheckBoxAll {
    // {
    //  all: selector for top checkbox
    //  one: selector for bottom checkbox
    //  change: callback on check
    // }
    constructor(props) {
        this.one = props.one;
        this.all = props.all;
        this.func = props.change;
        this.props = this.props;

        this.refresh();
    }

    refresh() {
        let your = this;
        let storage = [];

        // listen on all one.
        $(your.one).each(function (i, e) {
            $(e).prop('checked', false);
            $(e).on("change", function(e) {
                let data = $(this).attr('data');
                if ($(this).prop('checked')) {
                    storage.push(data);
                } else {
                    storage = storage.filter((v) => v !== data);
                }
                console.log(storage);
                your.func({store: storage});
            });
        });

        // listen on all and change to one.
        $(your.all).prop('checked', false);
        $(your.all).on("change", function(e) {
            storage = []; // empty
            if ($(this).prop('checked')) {
                $(your.one).each(function (i, e) {
                    //console.log($(element));
                    storage.push($(e).attr('data'));
                    $(e).prop('checked', true);
                });
            } else {
                $(your.one).each(function (i, e) {
                    $(e).prop('checked', false);
                });
            }
            console.log(storage);
            your.func({store: storage});
        });
    }
}