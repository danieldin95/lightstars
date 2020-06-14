
export class Base {
    // {
    //    id: ".container",
    //    default: "instances"
    //    force: false, // force to apply default.
    // }
    constructor(props) {
        this.id  = props.id;
        this.force = props.force;
        this.props = props;
        console.log('Base', props);
    }

    render() {
        $(this.id).html(this.template());
    }

    loading() {
    }

    template(v) {
        return (``)
    }

    child(id) {
        return this.id + " " + id;
    }

    title(name) {
        $(document).attr("title", name);
    }
}
