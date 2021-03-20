import {Ctl} from "./ctl.js"
import {LeasesCtl} from "./leases.js";
import {PortCtl} from "./port.js";

export class NetworkCtl extends Ctl {
    // {
    //   id: '#network'
    //   header: {
    //     id: '#'
    //  }
    //   leases: {
    //     id: '#leases'
    //   },
    // }
    constructor(props) {
        super(props);
        let name = $(this.id).attr("name");
        let uuid = $(this.id).attr("data");
        this.uuid = uuid;
        this.name = name;
        this.tasks = props.tasks || "tasks";

        this.leases = new LeasesCtl({...props.leases, uuid, name});
        this.ports = new PortCtl({...props.ports, name});
    }
}
