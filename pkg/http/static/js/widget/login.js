import {I18N} from "../lib/i18n.js";
import {Widget} from "./widget.js";
import {Location} from "../lib/location.js";


export class Login extends Widget {
    // {
    //   parent: "#container"
    // }
    constructor(props) {
        super(props);
        this.props = props;
        this.parent = props.parent;

        this.refresh();
    }

    compile(tmpl, data) {
        return template.compile(tmpl)(data);
    }

    refresh() {
        this.title(I18N.i('login'));
        $(this.parent).html(this.render({}));
        $(this.parent).find('#autofocus').focus();
    }

    next() {
        let page = Location.get();
        let query = Location.query();
        return page + "?" + query
    }

    render(v) {
        return this.compile(`
        <form action="" method="post">
        <input type="text" class="d-none" name="next" value="${this.next()}"/>
        <div class="card login">
            <div class="card-header">
                <div class="row">
                    <div class="col-sm-3">
                        <img class="" src="/static/images/lightstar.png" alt="">
                    </div>
                    <div class="col-sm-7">
                        <span style="color: #ffffff;">
                            <strong><a href="https://github.com/danieldin95/lightstar">{{'lightstars' | i}}</a></strong>
                        </span>
                        <br/>
                        <span style="color: #ced4da">
                            <small>{{'a simple and small platform to manage your cloud' | i}}</small>
                        </span>
                    </div>
                </div>
            </div>
            <div class="card-body">
                <div class="form-group row">
                    <label for="name" class="col-sm-3 col-form-label-sm">{{'username' | i}}</label>
                    <div class="col-sm-7">
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" name="name" value="" id="autofocus"/>
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
