import settings

from django.template.response import SimpleTemplateResponse
from django.views.decorators.csrf import ensure_csrf_cookie

@ensure_csrf_cookie
def render(request):
    env = 'debug' if settings.DEBUG else 'build'
    return SimpleTemplateResponse('../res/index.html', {'env': env})