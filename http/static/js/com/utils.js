
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

export function ListenChangeAll(storage, one, all, fn) {
    // listen on all one.
    for (let i = 0; i < $(one).length; i++) {
        $(one).eq(i).on("change", storage, function(e) {
            let data = $(this).attr("data");
            if ($(this).prop("checked")) {
                e.data.push(data)
            } else {
                e.data = e.data.filter(v => v != data);
            }
            fn({data: storage})
        });
    }

    // listen on all and change to one.
    $(all).on("change", storage, function(e) {
        if ($(this).prop("checked")) {
            $(one).each(function (index, element) {
                e.data.push($(this).attr("data"));
                $(element).prop("checked", true);
            });
        } else {
            $(one).each(function (index, element) {
                e.data = [];
                $(element).prop("checked", false);
            });
        }
        fn({data: storage})
    });
}