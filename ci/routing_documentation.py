from os import path, environ
import unittest
import subprocess
import json
from typing import TypedDict

INDIGO_BACKEND_EXECUTABLE_LOCATION: str = environ["INDIGO_BACKEND_EXECUTABLE_LOCATION"]
OPENAPI_SPEC_LOCATION: str = environ["OPENAPI_SPEC_LOCATION"]

def descend(root: dict, target: str):
    """
    Descend into a key of an unknown dict. Returns None if the key does not exist.
    :param root: The source dict
    :param target: The key to descend into
    """
    if not target in root.keys():
        return None
    return root[target]

def join_paths(parent: str, child: str) -> str:
    """
    Join two path strings where the child is relative to the parent
    """
    if parent != "":
        if parent[-1] == "/" and child[0] == "/":
            return parent + child[1:]
        elif parent[-1] != "/" and child[0] != "/":
            return parent + "/" + child
        else:
            return parent + child
    else:
        return child


class ImplementedRouterNode(TypedDict):
    path: str
    methods: list[str]
    children: list["ImplementedRouterNode"]


class ImplementedRouteIsDocumented(unittest.TestCase):
    implemented_routes: ImplementedRouterNode

    def setUp(self):
        assert INDIGO_BACKEND_EXECUTABLE_LOCATION, "INDIGO_BACKEND_EXECUTABLE_LOCATION is not defined"
        assert OPENAPI_SPEC_LOCATION, "OPENAPI_SPEC_LOCATION is not defined"

        # Get executable location
        executable_location = INDIGO_BACKEND_EXECUTABLE_LOCATION
        assert path.exists(executable_location), "specified executable does not exist"
        assert path.isfile(executable_location), "specified executable is not a file"
        executable_location_absolute = path.abspath(executable_location)

        # Get spec location
        spec_location = OPENAPI_SPEC_LOCATION
        assert path.exists(spec_location), "specified spec file does not exist"
        assert path.isfile(spec_location), "specified spec file is not a file"
        spec_location_absolute = path.abspath(spec_location)

        # Dump routes and capture output
        self.implemented_routes = json.loads(subprocess.check_output(
            args=f"{executable_location_absolute} dump_routes",
            shell=True,
            text=True
        ))

        # Read documented routes from fs
        with open(spec_location_absolute, 'r') as f:
            self.spec_routes = json.loads(f.read())

        # Grab openAPI version
        assert isinstance(descend(self.spec_routes, 'openapi'), str), "OpenAPI version not found"

        spec_paths = descend(self.spec_routes, "paths")
        assert spec_paths is not None, "paths key not found!"
        assert isinstance(spec_paths, dict), "paths value is not dict"

    def _recursive_check_documented(self, node: ImplementedRouterNode, parent_path: str):
        complete_path = join_paths(parent_path, node["path"])

        # Assert route is documented, assert methods match
        self.assertIn(
            member=complete_path,
            container=self.spec_routes["paths"].keys(),
            msg=f"Route {complete_path} does not exist in OpenAPI spec!"
        )
        for method in node["methods"]:
            self.assertIn(
                member=method.lower(),
                container=self.spec_routes["paths"][complete_path].keys(),
                msg=f"Route {complete_path} does not have documentation for {method} operation!"
            )
            operation = self.spec_routes["paths"][complete_path][method.lower()]
            self.assertIsInstance(
                obj=descend(operation, "summary"),
                cls=str,
                msg=f"Route {complete_path} does not have a summary!"
            )
            self.assertIsInstance(
                obj=descend(operation, "description"),
                cls=str,
                msg=f"Route {complete_path} does not have a description!"
            )
            self.assertIsInstance(
                obj=descend(operation, "responses"),
                cls=dict,
                msg=f"Route {complete_path} does not have documented responses!"
            )
            self.assertGreater(
                a=len(descend(operation, "responses").keys()),
                b=0,
                msg=f"Route {complete_path} has no documented responses!"
            )

        # Iterate over children if applicable
        if descend(node, "children"):
            for child in node["children"]:
                self._recursive_check_documented(child, complete_path)


    def test_all_routes_documented(self):
        """
        Test all routes defined in indigo_backend to see if they are defined in
        the OpenAPI spec. Requires that each method is documented for each path
        and requires that each method for each path has a summary, description, and
        return code list.
        """
        self._recursive_check_documented(self.implemented_routes, "/")

    def _get_implemented_methods_for_documented_route(self, route: str) -> list[str]|None:
        # Break down route into core components
        steps: list[str] = route.strip("/").split("/")
        if "" in steps:
            steps.remove("")
        steps.reverse()
        target = self.implemented_routes

        while len(steps) > 0:
            step = steps.pop()
            available_steps = [x["path"].replace("/", "") for x in target["children"]] if descend(target, "children") else []
            if step in available_steps:
                target = target["children"][available_steps.index(step)]
            else:
                return None

        return [x.lower() for x in target["methods"]]

    def test_all_routes_implemented(self):
        """
        Check that each route and method declared in the OpenAPI spec is implemented
        in indigo_backend.
        """
        for route, methods in self.spec_routes["paths"].items():
            implemented_methods = self._get_implemented_methods_for_documented_route(route)
            self.assertIsNotNone(
                obj=implemented_methods,
                msg=f"Route {route} from OpenAPI spec is not implemented in backend!"
            )
            for method in methods.keys():
                self.assertIn(
                    member=method,
                    container=implemented_methods,
                    msg=f"Method {method} is not implemented in backend for route {route} as dictated by OpenAPI spec!"
                )



if __name__ == "__main__":
    unittest.main()