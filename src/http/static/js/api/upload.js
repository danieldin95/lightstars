import {Api} from "./api.js";
import {Alert} from "../com/alert.js";
import {ProgressBar} from "../com/progressbar.js"

export class UploadApi extends Api {
    // {
    //   id: '#'
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
            return super.url(`/upload/${uuid}`);
        }
        return super.url(`/upload`);
    }

    //
    // data: new FormData($('form')[0])
    // event: {start: func, process: func}
    //
    upload(data) {
        let url = this.url(this.uuids[0]);
        let now = new Date();

        let event = {
            id: this.id,
            tasks: this.tasks,
            bar: now.getTime(),
            progress: function(e) {
                if (e.lengthComputable) {
                    $(this.id).find(`#${this.bar} .progress-bar`).css('width', `${e.loaded / e.total * 100}%`);
                }
            },
            success: function(e) {
                $(this.id).find(`#${this.bar}`).remove();
                $(this.tasks).append(Alert.success(`upload ${e.message}`));
            },
            start: function(e) {
                $(this.id).append(ProgressBar(this.bar));
            },
            fail: function(e) {
                $(this.id).find(`#${this.bar}`).remove();
                $(this.tasks).append(Alert.danger((`POST ${url}: ${e.responseText}`)));
            }
        };
        $.ajax({
            url: url,
            type: 'POST',
            data: data,
            contentType: false,
            processData: false,
            success: function (resp) { event.success(resp); },
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
