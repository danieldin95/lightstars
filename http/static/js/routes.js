import {Index} from "./widget/container/index.js";
import {Guest} from "./widget/container/guest.js";
import {Location} from "./com/location.js";


export class Routes {
    // {
    //   hyper: $,
    //   container: ''
    // }
    constructor(props) {
        this.routes = [
            {
                prefix: "/instance/",
                function: (p) => {
                    let uuid = p.split('/', 3)[2];
                    new Guest({
                        id: props.container,
                        uuid: uuid,
                    });
                },
            },
            {
                prefix: "/datastore/",
                function: (p) => {
                    //TODO
                },
            },
            {
                prefix: "/network/",
                function: (p) => {
                    //TODO
                },
            },
            {
                prefix: "",
                function: (p) => {
                    new Index({
                        id: props.container,
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
