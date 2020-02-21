
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