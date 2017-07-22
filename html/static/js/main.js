$(function($) {
    render();

    $('#save-btn').on('click', function(e) {
        e.preventDefault();

        var key = $('#k').val(),
            val = $('#v').val();

        if (key.length == 0) {
            notify("键不能为空", "danger");
            return;
        }

        if (val.length == 0) {
            notify("值不能为空", "danger");
            return;
        }

        $.ajax({
            url    : "/set",
            method : "POST",
            data   : {"key": key, "value": val},
            success: function() {
                render();
                $('#k').val('');
                $('#v').val('');
                notify("添加成功", "success");
            },
            error  : function(xhr) {
                notify("添加失败: " + xhr.responseText, "danger");
            }
        });
    });

    $('#update-btn').on('click', function(e) {
        e.preventDefault();

        var key = $('#k').val(),
            val = $('#v').val();

        if (key.length == 0) {
            notify("键不能为空", "danger");
            return;
        }

        if (val.length == 0) {
            notify("值不能为空", "danger");
            return;
        }

        $.ajax({
            url    : "/update",
            method : "POST",
            data   : {"key": key, "value": val},
            success: function() {
                render();
                $('#k').val('');
                $('#v').val('');
                notify("更新成功", "success");
            },
            error  : function(xhr) {
                notify("更新失败: " + xhr.responseText, "danger");
            }
        });
    });

    $('#del-btn').on('click', function(e) {
        e.preventDefault();

        var key = $('#k').val();

        if (key.length == 0) {
            notify("键不能为空", "danger");
            return;
        }

        $.ajax({
            url    : "/del",
            method : "POST",
            data   : {"key": key},
            success: function() {
                render();
                $('#k').val('');
                notify("删除成功", "success");
            },
            error  : function(xhr) {
                notify("删除失败: " + xhr.responseText, "danger");
            }
        });
    });

    function notify(body, level) {
        var elem = $('<div class="alert alert-' + level + '" role="alert">' + body + '<div>');
        $('.navbar').after(elem);
        elem.alert();
        setTimeout(function() {
            elem.alert('close');
        }, 1500);
    }

    function render() {
        $.getJSON('/all', function(data) {
            renderTable(data);
        }).fail(function() {
            notify("数据获取失败", "danger");
        });
    }

    function renderTable(data) {
        var elem = $('.easy-tree > ul');
        elem.empty();

        $(data).each(function(_, kv) {
            var k = kv[0], v = kv[1];
            renderTree(elem, k, v, 0);
        });

        elem.EasyTree();
    }

    function renderTree(root, k, v, depth) {
        var i = k.indexOf('.');
        if (i == -1) { // leaf node
            root.append('<li class="leaf">' + k + ': ' + v + '</li>');
        } else {
            var pk = k.slice(0, i),
                ck = k.slice(i+1),
                selector = 'li.' + pk,
                elem = root.find(selector);

            if (!elem.length) {
                elem = $('<li class="' + pk + '">' + pk + '<ul></ul></li>');
                root.append(elem);
            }

            renderTree(elem.find('ul'), ck, v, depth+1);
        }
    }
});
