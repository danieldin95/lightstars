import {Ctl} from "./ctl.js"
import {NetworkApi} from "../api/network.js";
import {LeasesCtl} from "./leases.js";

export class NetworkCtl extends Ctl {
    // {
    //   id: '#network'
    //   header: {
    //     id: '#'
    //  }
    //   leases: {
    //     id: '#leases'
    //   },
    //   subnets: {
    //     id: "#subnets"
    //   },
    // }
    constructor(props) {
        super(props);
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";

        console.log("NetworkCtl", this.props, $(this.id));
        this.leases = new LeasesCtl({id: props.leases.id, uuid, name});
    }
}
