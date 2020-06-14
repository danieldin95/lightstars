export class Ctl {
    // {
    //    id: "#xx"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;

        console.log("Ctl", props);
    }

    child(id) {
        return [this.id, id].join(" ")
    }
}
