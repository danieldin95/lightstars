
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
        $(one).eq(i).prop('checked', false);
        $(one).eq(i).on("change", function(e) {
            let data = $(this).attr('data');
            if ($(this).prop('checked')) {
                storage.push(data);
            } else {
                storage = storage.filter(v => v != data);
            }
            fn({data: storage});
        });
    }

    // listen on all and change to one.
    $(all).prop('checked', false);
    $(all).on("change", function(e) {
        storage.splice(0, storage.length); // empty
        if ($(this).prop('checked')) {
            $(one).each(function (index, element) {
                //console.log($(element));
                storage.push($(element).attr('data'));
                $(element).prop('checked', true);
            });
            //console.log(storage);
        } else {
            $(one).each(function (index, element) {
                $(element).prop('checked', false);
            });
        }
        fn({data: storage});
    });
}