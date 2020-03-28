
export class WidgetBase {
    constructor(props) {
        this.id = props.id;
        this.props = props;
        this.tasks = props.tasks;
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }
}