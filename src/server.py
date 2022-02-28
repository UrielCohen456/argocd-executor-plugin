import json
from http.server import BaseHTTPRequestHandler, HTTPServer


class Plugin(BaseHTTPRequestHandler):

    def args(self):
        return json.loads(self.rfile.read(int(self.headers.get('Content-Length'))))

    def reply(self, reply):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(json.dumps(reply).encode("UTF-8"))

    def error(self, message: str):
        self.reply({
            'node': {
                'phase': 'Error', 
                'message': str.format("ArgoCD executor plugin error: {}", message),
            }
        })

    def success(self, message: str, outputs: list = None):
        ret_val = {
            'node': {
                'phase': 'Succeeded', 
                'message': str.format("ArgoCD executor plugin succeeded: {}", message),
            }
        }

        if outputs is not None:
            ret_val['node']["outputs"] = { 'parameters': outputs }

        self.reply(ret_val)

    def fail(self, message: str):
        self.reply({
            'node': {
                'phase': 'Failed', 
                'message': str.format("ArgoCD executor plugin failed: {}", message),
            }
        })

    def do_POST(self):
        if self.path != '/api/v1/template.execute':
            self.send_response(404)
            self.end_headers()
            return

        args = self.args()

        if 'argocd' not in args['template'].get('plugin', {}):
            self.reply({})
            return

        self.success("synced app")


if __name__ == '__main__':
    httpd = HTTPServer(('', 4355), Plugin)
    httpd.serve_forever()
