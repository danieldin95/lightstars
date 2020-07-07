export class Login {
    // {
    //   parent: "#container"
    // }
    constructor(props) {
        this.props = props;
        this.parent = props.parent;

        this.refresh();
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }

    refresh() {
        console.log(this.parent, $(this.parent));
        $(this.parent).html(this.render({}));
    }

    render(v) {
        return this.compile(`
        <form action="/ui/login" method="post">
        <div class="card login">
            <div class="card-header">
                <div class="row">
                    <div class="col-sm-3">
                        <img class="" src="/static/images/lightstar-6.png" alt="">
                    </div>
                    <div class="col-sm-7">
                        <span>
                            <strong><a href="https://github.com/danieldin95/lightstar">LightStar</a></strong>
                        </span>
                        <br/>
                        <span>
                            <small style="color: #ced4da">A simple and small platform to manage your cloud</small>
                        </span>
                    </div>
                </div>
            </div>
            <div class="card-body">
                <div class="form-group row">
                    <label for="name" class="col-sm-3 col-form-label-sm">{{'username' | i}}</label>
                    <div class="col-sm-7">
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value="" autofocus/>
                        </div>
                    </div>
                </div>
                <div class="form-group row">
                    <label for="password" class="col-sm-3 col-form-label-sm">{{'password' | i}}</label>
                    <div class="col-sm-7">
                        <div class="input-group">
                            <input type="password" class="form-control form-control-sm" name="password" value=""/>
                        </div>
                    </div>
                </div>
            </div>
            <div class="card-footer d-flex justify-content-end">
                <button id="reset" class="btn btn-outline-dark btn-sm mr-2" type="reset">{{'reset' | i}}</button>
                <button id="submit" class="btn btn-outline-success btn-sm mr-0">{{ 'login' | i}}</button>
            </div>
        </div>
        </form>`, v)
    }
}
