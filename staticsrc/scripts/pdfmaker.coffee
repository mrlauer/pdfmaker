require.config
    paths :
        jquery : 'lib/jquery-1.6.1.min'
        backbone : 'lib/backbone'
        underscore : 'lib/underscore'
    priority : [ 'jquery', 'underscore', 'backbone']

#syntax looks funny, i know
require [ 'jquery', 'order!underscore', 'order!backbone' ], -> $ ->
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

            @render()

        render: ->

        changeText: => @model.save 'Text', $('#text').val()

        events: { 'change #text' : 'changeText' }

    model = new Document
    model.fetch()
    doc_view = new DocView { model: model }

