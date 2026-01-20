from django.urls import include, path
import views.index.view as index_view
import controllers.urls as urls

urlpatterns = [
    path('', index_view.render),
    path('api/', include(urls)),
]
