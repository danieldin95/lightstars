import {Api} from "./api.js";
import {Alert} from "../com/alert.js";
import {ProgressBar} from "../com/progressbar.js"

export class UploadApi extends Api {
    // {
    //  id: '#'
    //   uuid: [datastore],
    //   tasks: 'tasks',
    //   name: ''
    // }
    constructor(props) {
        super(props);
        this.id = props.id;
    }

    url(uuid) {
        if (uuid) {
            return `/api/upload/${uuid}`
        }
        return `/api/upload`
    }

    //
    // data: new FormData($('form')[0])
    // event: {start: func, process: func}
    //
    upload(data) {
        let event = {
            progress: (e) => {
                if (e.lengthComputable) {
                    $(`${this.id} .progress-bar`).css('width', `${e.loaded / e.total * 100}%`);
                }
            },
            success: (e) => {
                $(this.id).empty();
                $(this.tasks).append(Alert.success(`upload iso ${data.message}`));
            },
            start: (e) => {
                $(this.id).html(ProgressBar());
            },
            fail: (e) => {
                $(this.id).empty();
                $(this.tasks).append(Alert.danger((`POST upload file: ${e.responseText}`)));
            }
        };
        $.ajax({
            url: this.url(this.uuids[0]),
            type: 'POST',
            data: data,
            contentType: false,
            processData: false,
            success: function (file) { event.success(file); },
            xhr: function () {
                event.start();

                let x = $.ajaxSettings.xhr();
                x.upload.addEventListener("progress", function (e) {
                    event.progress(e);
                }, false);

                return x;
            },
        }).fail(function (e) { event.fail(e); });
    }
}