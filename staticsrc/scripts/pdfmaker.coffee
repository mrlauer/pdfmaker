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

    propertyNames =
        Font: 'Font'
        FontSize: 'Font Size'
        BaselineSkip: 'Baseline Skip'
        LeftMargin: 'Left Margin'
        RightMargin: 'Right Margin'
        TopMargin: 'Top Margin'
        BottomMargin: 'Bottom Margin'
        PageWidth: 'Page Width'
        PageHeight: 'Page Height'
        Text: 'Text'

    sizeControlFields =
        FontSize : true
        BaselineSkip : true
        LeftMargin : true
        RightMargin : true
        TopMargin : true
        BottomMargin : true
        PageWidth : true
        PageHeight : true

    sizeControls = ({ name: name, label: propertyNames[name] } for name in [
        'FontSize'
        'BaselineSkip'
        'LeftMargin'
        'RightMargin'
        'TopMargin'
        'BottomMargin'
        'PageWidth'
        'PageHeight'
    ])

    class Document extends Backbone.Model
        initialize: (args) ->
            @id = args?.id

        defaults: defaultDoc

        urlRoot: -> '/document/'

        validate: (attrs) ->
            sizeRE = /^\s*(\d+(\.\d*)?|\.\d+)\s*("|in|pt)\s*$/
            for field, val of attrs
                if field of sizeControlFields
                    if val? and !sizeRE.test val
                        return "Bad value for #{field}"
            return null


    class DocView extends Backbone.View
        el: $ 'body'

        initialize: (args) ->
            _.bindAll @

            @model = (args?.model) ? new Document
            @model.on 'change', -> doc_view.render()

            @render()

        render: ->
            model = @model
            templ = mustache.render doctempl,
                fonts: availableFonts
                sizeControls: sizeControls
                get: -> (key, render)-> _.escape model.get render key
            $('#content-div').html templ
            $('#Font').val @model.get 'Font'
            @

        changeText: => @model.save 'Text', $('#Text').val()
        changeProp: (prop) =>
            self = @
            attrs = {}
            attrs[prop] = @$("##{prop}").val()
            @model.save attrs,
                { error: -> self.$("##{prop}").addClass 'error' }

        events:
            'change #Text' : 'changeText'
            'change .docControl' : (ev) -> @changeProp $(ev.currentTarget).attr('name')

    doc_view = new DocView

