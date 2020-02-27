

export function Arr2Json(data) {
    let ret = {};

    data.forEach(function (element, index) {
        let key = element['name'];
        let value = element['value'];
        if (key && value) {
            ret[key] = value;
        }
    });
    return ret;
}


export function ListenChangeAll(one, all, fn) {
    let storage = [];

    // listen on all one.
    for (let i = 0; i < $(one).length; i++) {
        $(one).eq(i).prop('checked', false);
        $(one).eq(i).on("change", function(e) {
            let data = $(this).attr('data');
            if ($(this).prop('checked')) {
                storage.push(data);
            } else {
                storage = storage.filter((v) => v !== data);
            }
            fn({store: storage});
        });
    }

    // listen on all and change to one.
    $(all).prop('checked', false);
    $(all).on("change", function(e) {
        storage = []; // empty
        if ($(this).prop('checked')) {
            $(one).each(function (i, e) {
                //console.log($(element));
                storage.push($(e).attr('data'));
                $(e).prop('checked', true);
            });
            //console.log(storage);
        } else {
            $(one).each(function (i, e) {
                $(e).prop('checked', false);
            });
        }
        fn({store: storage});
    });
}

export function PrintN(num, n) {
    return (Array(n).join(0) + num).slice(-n);
}