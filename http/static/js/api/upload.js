import {Api} from "./api.js";
import {AlertDanger, AlertSuccess} from "../com/alert.js";
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
        let your = this;

        let event = {
            id: this.id,
            progress: function (e) {
                if (e.lengthComputable) {
                    $(`${this.id} .progress-bar`).css('width', `${e.loaded / e.total * 100}%`);
                }
            },
            success: function (e) {
                $(this.id).empty();
                $(this.tasks).append(AlertSuccess(`upload iso ${data.message}`));
            },
            start: function (e) {
                $(this.id).html(ProgressBar());
            },
            fail: function (e) {
                $(this.id).empty();
                $(your.tasks).append(AlertDanger((`POST upload file: ${e.responseText}`)));
            }
        };
        $.ajax({
            url: your.url(your.uuids[0]),
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