
export let Utils = {
    toN: function (num, n) {
        return (Array(n).join(0) + num).slice(-n);
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


export class CheckBoxTop {
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
        for (let i = 0; i < $(your.one).length; i++) {
            $(your.one).eq(i).prop('checked', false);
            $(your.one).eq(i).on("change", function(e) {
                let data = $(this).attr('data');
                if ($(this).prop('checked')) {
                    storage.push(data);
                } else {
                    storage = storage.filter((v) => v !== data);
                }
                your.func({store: storage});
            });
        }

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
                //console.log(storage);
            } else {
                $(your.one).each(function (i, e) {
                    $(e).prop('checked', false);
                });
            }
            your.func({store: storage});
        });
    }
}