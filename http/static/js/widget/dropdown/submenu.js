
export class SubMenu {
    //
    constructor(props) {
        $('.dropdown-submenu > a').on("click", function(e) {
            $('.dropdown-submenu .dropdown-menu').removeClass('show');
            $(this).next('.dropdown-menu').addClass('show');
            e.stopPropagation();
        });

        $('.dropdown').on("hidden.bs.dropdown", function() {
            // hide any open menus when parent closes
            $('.dropdown-menu.show').removeClass('show');
        });
    }
}

