require.config
    paths :
        jquery : 'lib/jquery-1.6.1.min'
        backbone : 'lib/backbone'
        underscore : 'lib/underscore'
        mustache : 'lib/requirejs.mustache'
    priority : [ 'jquery', 'underscore', 'backbone']

#syntax looks funny, i know
require [ 'mustache', 'order!jquery', 'order!underscore', 'order!backbone' ],
  (mustache) -> $ ->
    
    docTempl = """<textarea id="text" name="text">{{text}}</textarea>"""
        
    doc_view = null

    class Document extends Backbone.Model
        initialize: (args) ->
            @id = args?.id

        defaults:
            LeftMargin: 1
            TopMargin: 1
            Font: "Times New Roman"
            Text: """Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."""

        urlRoot: -> '/document/'


    class DocView extends Backbone.View
        el: $ 'body'

        initialize: (args) ->
            _.bindAll @

            @model = (args?.model) ? new Document
            @model.on 'change', -> doc_view.render()

            @render()

        render: ->
            templ = mustache.render docTempl,
                text: @model.get 'Text'
            $('#content-div').html templ

        changeText: => @model.save 'Text', $('#text').val()

        events: { 'change #text' : 'changeText' }

    model = new Document
    doc_view = new DocView { model: model }
    model.fetch()

