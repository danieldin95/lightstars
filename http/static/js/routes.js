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
                prefix: "instance/",
                function: (p) => {
                    let uuid = p.split('/', 2)[1];
                    new Guest({
                        id: props.container,
                        uuid: uuid,
                    });
                },
            },
            {
                prefix: "datastore/",
                function: (p) => {
                    //TODO
                },
            },
            {
                prefix: "network/",
                function: (p) => {
                    //TODO
                },
            },
            {
                prefix: "",
                function: (p) => {
                    new Index({
                        id: props.container,
                        default: "instances",
                        name: props.hyper.attr('name'),
                    });
                },
            },
        ];
        this.render();
    }

    render() {
        let cur = Location.get();
        for (let i = 0; i < this.routes.length; i++) {
            let rte = this.routes[i];
            if (cur.startsWith(rte.prefix)) {
                rte.function(cur);
                break;
            }
        }
    }
}