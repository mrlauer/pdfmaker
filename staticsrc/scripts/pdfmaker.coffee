require.config
    paths :
        jquery : 'lib/jquery-1.6.1.min'
        backbone : 'lib/backbone'
        underscore : 'lib/underscore'
        mustache : 'lib/requirejs.mustache'
    priority : [ 'jquery', 'underscore', 'backbone']

#syntax looks funny, i know
require [ 'mustache', 'text!doctempl.html', 'order!jquery', 'order!underscore', 'order!backbone' ],
  (mustache, doctempl) -> $ ->
    
    doc_view = null

    class Document extends Backbone.Model
        initialize: (args) ->
            @id = args?.id

        urlRoot: -> '/document/'


    class DocView extends Backbone.View
        el: $ 'body'

        initialize: (args) ->
            _.bindAll @

            @model = (args?.model) ? new Document
            @model.on 'change', -> doc_view.render()

            @render()

        render: ->
            templ = mustache.render doctempl,
                attr: @model.attributes
            $('#content-div').html templ

        changeText: => @model.save 'Text', $('#text').val()

        events: { 'change #text' : 'changeText' }

    model = new Document
    doc_view = new DocView { model: model }
    model.fetch()

