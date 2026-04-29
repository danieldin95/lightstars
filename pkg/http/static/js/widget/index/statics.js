import {widget} from "../widget.js";
import {HyperApi} from "../../api/hyperapi.js";


export class statics extends widget {
    // {
    //   id: '#xx'.
    // }
    constructor(props) {
        super(props);
        console.log(props);
    }

    refresh(data, func) {
        if (typeof data == 'function') {
            func = data;
            data = {};
        }
        new HyperApi({tasks: this.tasks}).statics(this,function (e) {
            $(e.data.id).html(e.data.render(e.resp));
            if (func) {
                func({data, resp: e.resp});
            }
        });
    }

    render(data) {
        return this.compile(`
            <div class="dashboard-stats home-stats-row">
                <div class="dashboard-group home-stats-item">
                    <div class="dashboard-group-title">{{'datastore' | i}}</div>
                    <div class="dashboard-grid">
                        <div class="dashboard-stat total">
                            <div class="label">total</div>
                            <div class="value">{{datastore.total}}</div>
                        </div>
                        <div class="dashboard-stat up">
                            <div class="label">active</div>
                            <div class="value">{{datastore.active}}</div>
                        </div>
                    </div>
                </div>
                <div class="dashboard-group home-stats-item">
                    <div class="dashboard-group-title">{{'guest instances' | i}}</div>
                    <div class="dashboard-grid home-stat-grid-two-row">
                        <div class="dashboard-stat total">
                            <div class="label">total</div>
                            <div class="value">{{instance.total}}</div>
                        </div>
                        <div class="dashboard-stat up">
                            <div class="label">running</div>
                            <div class="value">{{instance.active}}</div>
                        </div>
                        <div class="dashboard-stat down">
                            <div class="label">shutdown</div>
                            <div class="value">{{instance.inactive}}</div>
                        </div>
                    </div>
                </div>
                <div class="dashboard-group home-stats-item">
                    <div class="dashboard-group-title">{{'virtual networks' | i}}</div>
                    <div class="dashboard-grid">
                        <div class="dashboard-stat total">
                            <div class="label">total</div>
                            <div class="value">{{network.total}}</div>
                        </div>
                        <div class="dashboard-stat up">
                            <div class="label">active</div>
                            <div class="value">{{network.active}}</div>
                        </div>
                    </div>
                </div>
                <div class="dashboard-group home-stats-item">
                    <div class="dashboard-group-title">{{'virtual ports' | i}}</div>
                    <div class="dashboard-grid home-stat-grid-two-row">
                        <div class="dashboard-stat total">
                            <div class="label">total</div>
                            <div class="value">{{ports.total}}</div>
                        </div>
                        <div class="dashboard-stat up">
                            <div class="label">up</div>
                            <div class="value">{{ports.active}}</div>
                        </div>
                        <div class="dashboard-stat down">
                            <div class="label">down</div>
                            <div class="value">{{ports.inactive}}</div>
                        </div>
                    </div>
                </div>
            </div>`, data);
    }
}
