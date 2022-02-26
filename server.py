import json
from http.server import BaseHTTPRequestHandler, HTTPServer


class Plugin(BaseHTTPRequestHandler):

    def args(self):
        return json.loads(self.rfile.read(int(self.headers.get('Content-Length'))))

    def reply(self, reply):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps(reply).encode("UTF-8"))

    def unsupported(self):
        self.send_response(404)
        self.end_headers()

    def do_POST(self):
        if self.path != '/api/v1/template.execute':
            self.unsupported()
            return

        args = self.args()

        if 'argocd' not in args['template'].get('plugin', {}):
            self.reply({})
            return

        self.reply({
            'node': {
                'phase': 'Succeeded', 
                'message': 'Hello template!',
                'outputs': {
                    'parameters': [
                        {'name': 'foo', 'value': 'bar'}
                    ]
                }
            }
        })


if __name__ == '__main__':
    httpd = HTTPServer(('', 4355), Plugin)
    httpd.serve_forever()
