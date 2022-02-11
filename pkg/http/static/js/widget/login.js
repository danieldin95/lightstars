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
        this.user = $.cookie('user');
        if (!this.user) {
            this.user = ""
        }
        this.refresh();
        this.loading();
    }

    refresh() {
        this.title(I18N.i('login'));
        $(this.parent).html(this.render({}));
    }

    loading() {
        let parent = $(this.parent);
        parent.find('#name').focus();
        parent.find('#password-app').on('click',  function (e) {
            let current = parent.find('#password').attr('type');
            if ( current === 'password') {
                parent.find('#password-app i').addClass('bi-eye-slash');
                parent.find('#password-app i').removeClass('bi-eye');
                parent.find('#password').attr('type', 'text');
            } else {
                parent.find('#password-app i').addClass('bi-eye');
                parent.find('#password-app i').removeClass('bi-eye-slash');
                parent.find('#password').attr('type', 'password');
            }
            parent.find('#password').focus();
        });
        parent.find('#name-app').on('click',  function (e) {
            parent.find('#name').focus();
        });
    }

    next() {
        let page = Location.get();
        let query = Location.query();
        return page + "?" + query;
    }

    render(v) {
        return this.compile(`
        <form action="" method="post">
        <input type="text" class="d-none" name="next" value="${this.next()}"/>
        <div class="card login">
            <div class="card-header">
                <div class="row">
                    <div class="col-12">
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
                    <label for="name" class="col-12 col-form-label-sm">{{'username' | i}}</label>
                    <div class="col-12">
                        <div class="input-group">
                            <input type="text" class="form-control form-control-sm" id="name" name="name" value="${this.user}"/>
                            <div class="input-group-append">
                                <a href="javascript:void(0)" class="input-group-text input-group-sm" id="name-app">
                                    <i class="bi bi-pencil"></i>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="form-group row">
                    <label for="password" class="col-12 col-form-label-sm">{{'password' | i}}</label>
                    <div class="col-12">
                        <div class="input-group">
                            <input type="password" class="form-control form-control-sm" id="password" name="password" value=""/>
                            <div class="input-group-append">
                                <a href="javascript:void(0)" class="input-group-text input-group-sm" id="password-app">
                                    <i class="bi bi-eye"></i>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="card-footer d-flex justify-content-end">
                <button id="submit" class="btn btn-outline-success btn-sm mr-0">{{ 'login' | i}}</button>
            </div>
        </div>
        </form>`, v)
    }
}
