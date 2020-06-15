import {Home} from "./widget/container/home.js";
import {Guest} from "./widget/container/guest.js";
import {Network} from "./widget/container/network.js";
import {Location} from "./com/location.js";
import {Pool} from "./widget/container/pool.js";

export class Routes {
    // {
    //   hyper: $,
    //   container: ''
    // }
    constructor(props) {
        this.routes = [
            {
                prefix: "/guest/",
                function: (p) => {
                    let uuid = p.split('/', 3)[2];
                    new Guest({
                        parent: props.container,
                        uuid: uuid,
                    });
                },
            },
            {
                prefix: "/datastore/",
                function: (p) => {
                    let uuid = p.split('/', 3)[2];
                    new Pool({
                        parent: props.container,
                        uuid: uuid,
                    });
                },
            },
            {
                prefix: "/network/",
                function: (p) => {
                    let uuid = p.split('/', 3)[2];
                    new Network({
                        parent: props.container,
                        uuid: uuid,
                    });
                },
            },
            {
                prefix: "",
                function: (p) => {
                    new Home({
                        parent: props.container,
                        default: "/instances",
                        name: props.hyper.attr('name'),
                    });
                },
            },
        ];
        this.render();
    }

    render() {
        let cur = Location.get();
        for (let rte of this.routes) {
            if (cur.startsWith(rte.prefix)) {
                rte.function(cur);
                break;
            }
        }
    }
}
