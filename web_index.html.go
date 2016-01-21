package main

const (
	INDEX_HTML = `<!DOCTYPE html>
<html>
<head lang="en">
  <title>Request Baskets</title>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css">
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap-theme.min.css">
  <script src="https://code.jquery.com/jquery-2.1.4.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>

  <style>
    html { position: relative; min-height: 100%; }
    body { padding-top: 70px; margin-bottom: 60px; }
    .footer { position: absolute; bottom: 0; width: 100%; height: 60px; background-color: #f5f5f5; }
    .container .text-muted { margin: 20px 0; }
    h1 { margin-top: 2px; }
    #more { margin-left: 40px; }
  </style>

  <script>
  (function($) {
    var basketsCount = 0;

    function onAjaxError(jqXHR) {
      $("#error_message_label").html("HTTP " + jqXHR.status + " - " + jqXHR.statusText);
      $("#error_message_text").html(jqXHR.responseText);
      $("#error_message").modal();
    }

    function clearBaskets() {
      $("#baskets").html("");
      basketsCount = 0;
    }

    function addBaskets(data) {
      if (data && data.names) {
        var baskets = $("#baskets");
        var index, name;
        for (index = 0; index < data.names.length; ++index) {
          name = data.names[index];
          baskets.append("<li><a href='/web/" + name + "'>" + name + "</a></li>");
          basketsCount++;
        }

        if (data.has_more) {
          $("#more").removeClass("hide");
        } else {
          $("#more").addClass("hide");
        }
      }
    }

    function fetchBaskets() {
      $.get("/baskets?skip=" + basketsCount, function(data) {
        addBaskets(data);
      }).fail(onAjaxError);
    }

    function createBasket() {
      var basket = $("#basket_name").val();
      $.post("/baskets/" + basket, function(data) {
        sessionStorage.setItem("token_" + basket, data.token);
        $("#created_message_text").html("<p>Basket '" + basket +
          "' is successfully created!</p><p>Your token is: <mark>" + data.token + "</mark></p>");
        $("#basket_link").attr("href", "/web/" + basket);
        $("#created_message").modal();

        // refresh
        clearBaskets();
        fetchBaskets();
      }).fail(onAjaxError).always(function() {
        $("#basket_name").val("");
      });
    }

    function saveMasterToken() {
      var token = $("#master_token").val();
      $("#master_token").val("");
      $("#master_token_dialog").modal("hide");
      if (token) {
        sessionStorage.setItem("master_token", token);
      } else {
        sessionStorage.removeItem("master_token");
      }
      indicateMasterToken();
    }

    function indicateMasterToken() {
      var token = $("#token");
      if (sessionStorage.getItem("master_token")) {
        token.removeClass("btn-warning");
        token.addClass("btn-success");
      } else {
        token.removeClass("btn-success");
        token.addClass("btn-warning");
      }
    }

    // Initialization
    $(document).ready(function() {
      $("#create_basket").on("submit", function(event) {
        createBasket();
        event.preventDefault();
      });
      $("#refresh").on("click", function(event) {
        clearBaskets();
        fetchBaskets();
      });
      $("#token").on("click", function(event) {
        $("#master_token_dialog").modal();
      });
      $("#save_master_token").on("click", function(event) {
        saveMasterToken();
      });
      $("#fetch_more").on("click", function(event) {
        fetchBaskets();
      });

      indicateMasterToken();
      fetchBaskets();
    });
  })(jQuery);
  </script>
</head>
<body>
  <!-- Fixed navbar -->
  <nav class="navbar navbar-default navbar-fixed-top">
    <div class="container">
      <div class="navbar-header">
        <a id="refresh" class="navbar-brand" href="#">Request Baskets</a>
      </div>
      <div class="collapse navbar-collapse">
        <form class="navbar-form navbar-right">
          <button id="token" type="button" title="Master Token" class="btn btn-warning">
            <span class="glyphicon glyphicon-lock"></span>
          </button>
        </form>
      </div>
    </div>
  </nav>

  <!-- Error message -->
  <div class="modal fade" id="error_message" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content panel-danger">
        <div class="modal-header panel-heading">
          <button type="button" class="close" data-dismiss="modal">&times;</button>
          <h4 class="modal-title" id="error_message_label">HTTP error</h4>
        </div>
        <div class="modal-body">
          <p id="error_message_text"></p>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
        </div>
      </div>
    </div>
  </div>

  <!-- Created message -->
  <div class="modal fade" id="created_message" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content panel-success">
        <div class="modal-header panel-heading">
          <button type="button" class="close" data-dismiss="modal">&times;</button>
          <h4 class="modal-title" id="created_message_label">Created</h4>
        </div>
        <div class="modal-body" id="created_message_text">
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
          <a id="basket_link" class="btn btn-primary">Open Basket</a>
        </div>
      </div>
    </div>
  </div>

  <!-- Master token dialog -->
  <div class="modal fade" id="master_token_dialog" tabindex="-1">
    <div class="modal-dialog">
      <div class="modal-content panel-warning">
        <div class="modal-header panel-heading">
          <button type="button" class="close" data-dismiss="modal">&times;</button>
          <h4 class="modal-title">Master Token</h4>
        </div>
        <form>
        <div class="modal-body">
          <p>By providing the master token you will gain additional privileges.</p>
          <div class="form-group">
            <label for="master_token" class="control-label">Token:</label>
            <input type="password" class="form-control" id="master_token">
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-default" data-dismiss="modal">Cancel</button>
          <button type="button" class="btn btn-success" id="save_master_token">Authorize</button>
        </div>
        </form>
      </div>
    </div>
  </div>

  <div class="container">
    <div class="row">
      <div class="col-md-4">
        <h1>Baskets</h1>
      </div>
      <div class="col-md-5 col-md-offset-3">
        <form id="create_basket" class="navbar-form">
          <div class="form-group">
            <input id="basket_name" type="text" placeholder="Basket Name" class="form-control">
          </div>
          <button type="submit" class="btn btn-primary">Create</button>
        </form>
      </div>
    </div>
    <hr/>
    <div class="row">
      <ul id="baskets">
      </ul>
      <div id="more" class="hide">
        <a id="fetch_more" class="btn btn-default btn-xs">more...</a>
      </div>
    </div>
  </div>

  <footer class="footer">
    <div class="container">
      <p class="text-muted"><small>Powered by <a href="https://github.com/darklynx/request-baskets">request-baskets</a></small></p>
    </div>
  </footer>
</body>
</html>`
)
