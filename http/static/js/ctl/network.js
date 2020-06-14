import {NetworkApi} from "../api/network.js";
import {LeasesCtl} from "./leases.js";

export class NetworkCtl {
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
        this.id = props.id;
        this.props = props;
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "Tasks";

        console.log("NetworkCtl", this.props, $(this.id));

        this.leases = new LeasesCtl({id: props.leases.id, uuid, name});
    }

    edit(data) {
        new NetworkApi({uuids: this.uuid}).edit(data);
    }
}
