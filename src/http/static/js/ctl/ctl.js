export class Ctl {
    // {
    //    id: "#xx"
    // }
    constructor(props) {
        this.id = props.id;
        this.props = props;
        console.log("props", this.props);
    }

    child(id) {
        return [this.id, id].join(" ")
    }
}
