import json
from http.server import BaseHTTPRequestHandler, HTTPServer
from pprint import pprint
from kubernetes import client, config, watch

class Plugin(BaseHTTPRequestHandler):
    def args(self) -> any:
        """Returns the plugin template provided from the workflow"""

        return json.loads(self.rfile.read(int(self.headers.get('Content-Length'))))

    def reply(self, reply: dict, code: int) -> None:
        """Sends a reply back to the caller"""
        self.send_response(code)
        self.end_headers()
        self.wfile.write(json.dumps(reply).encode("UTF-8"))

    def success(self, message: str, outputs: list = None) -> None:
        """Sends a success message to the caller"""
        ret_val = {
            'node': {
                'phase': 'Succeeded', 
                'message': str.format("ArgoCD executor plugin succeeded: {}", message),
            }
        }

        if outputs is not None:
            ret_val['node']["outputs"] = { 'parameters': outputs }

        self.reply(ret_val, 200)

    def fail(self, message: str) -> None:
        """Sends a fail message to the caller"""
        self.reply({
            'node': {
                'phase': 'Failed', 
                'message': str.format("ArgoCD executor plugin failed: {}", message),
            }
        }, 503)

    def error(self, message: str) -> None:
        """Sends an error message to the caller"""
        self.reply({
            'node': {
                'phase': 'Error', 
                'message': str.format("ArgoCD executor plugin error: {}", message),
            }
        }, 400)

    def do_POST(self) -> None:
        """Receives the request from the caller and check if it's valid. If it is, executes the plugin"""
        if self.path != '/api/v1/template.execute':
            self.reply({}, 404)
            return

        args = self.args()
        pprint(args)

        if 'argocd' not in args['template'].get('plugin', {}):
            self.reply({})
            return

        try:
            self.execute(args)
        except Exception as e:
            self.error(str(e))


    def execute(self, args: dict) -> None:
        """Executes the plugin and talks to the argocd server"""
        config.load_incluster_config()
        v1 = client.CoreV1Api()
        current_namespace = open("/var/run/secrets/kubernetes.io/serviceaccount/namespace").read()

        ret = v1.list_namespaced_pod(current_namespace)
        for i in ret.items:
            print("%s\t%s\t%s" %
                (i.status.pod_ip, i.metadata.namespace, i.metadata.name))

        # argocd = ("/usr/local/bin/argocd", "-h")
        # popen = subprocess.Popen(argocd, stdout=subprocess.PIPE)
        # popen.wait()
        # output = popen.stdout.read()

        self.success("synced app")


if __name__ == '__main__':
    httpd = HTTPServer(('', 4355), Plugin)
    httpd.serve_forever()
