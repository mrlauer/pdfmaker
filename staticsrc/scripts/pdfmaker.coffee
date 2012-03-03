require.config
    paths :
        jquery : 'lib/jquery-1.6.1.min'
        jqueryui : 'lib/jquery-ui-1.8.13.custom.min'
        backbone : 'lib/backbone'
        underscore : 'lib/underscore'
        mustache : 'lib/requirejs.mustache'
    priority : [ 'jquery', 'underscore', 'backbone']

#syntax looks funny, i know
require [ 'mustache', 'text!doctempl.html', 'order!jquery', 'order!jqueryui',
        'order!underscore', 'order!backbone' ],
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

        pdfUrl: ->
            if @id?
                "/pdf/#{@id}"
            else
                "/pdf/"

        validate: (attrs) ->
            for field, val of attrs
                if field of sizeControlFields
                    if val? and !lengthRE.test val
                        return "Bad value for #{field}"
            return null


    class DocView extends Backbone.View
        el: $ 'body'

        initialize: (args) ->
            _.bindAll @

            @model = (args?.model) ? new Document
            self = @
            @model.on 'change', ->
                self.render()
                if ! self.model.isNew()
                    router.navigate "edit/#{self.model.id}"

            @render()

        render: ->
            model = @model
            templ = mustache.render doctempl,
                fonts: availableFonts
                sizeControls: sizeControls
                get: -> (key, render)-> _.escape model.get render key
            @$('#content-div').html templ
            @$('#getPdf').button()
            @$('#Font').val @model.get 'Font'
            @

        changeText: => @model.save 'Text', $('#Text').val()
        changeProp: (prop) =>
            self = @
            attrs = {}
            attrs[prop] = @$("##{prop}").val()
            @model.save attrs,
                { error: -> self.$("##{prop}").addClass 'error' }

        getPdf: ->
            model = @model
            model.save {},
                success : ->
                    url = model.pdfUrl()
                    window.location = url

        events:
            'change #Text' : 'changeText'
            'change .docControl' : (ev) -> @changeProp $(ev.currentTarget).attr('name')
            'click  #getPdf' : 'getPdf'

    doc_view = new DocView

    class DocRouter extends Backbone.Router
        routes:
            "edit/:id": "edit"

        edit: (idstr)->
            id = parseInt idstr
            doc_view.model = new Document
                id : id
            doc_view.model.fetch
                success: -> doc_view.render()

    router = new DocRouter

    Backbone.history.start( { pushState: true} )



