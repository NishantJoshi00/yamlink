from locust import HttpUser, task, between


class ApiUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        with open("mapping.yaml", "w") as f:
            f.write("example: https://www.example.com")

    # def on_start2(self):
    #     with open("mapping.yaml", "w") as f:
    #         f.write("example: https://google.com")

    @task
    def my_task(self):
        self.first_task()

    def first_task(self):
        with self.client.get(
            "/example", allow_redirects=False, catch_response=True
        ) as output:
            if output.status_code == 302:
                redirect_url = output.headers["Location"].strip()
                if redirect_url != "https://www.example.com":
                    output.failure(
                        f"Expected redirect to https://www.example.com but got {redirect_url}"
                    )
            else:
                output.failure(f"Expected a redirect but got {output.status_code}")

    # def second_task(self):
    #     with self.client.get(
    #         "/example", allow_redirects=False, catch_response=True
    #     ) as output:
    #         if output.status_code == 302:
    #             redirect_url = output.headers["Location"]
    #             if redirect_url == "https://google.com":
    #                 output.failure(
    #                     f"Expected redirect to https://google.com but got {redirect_url}"
    #                 )
    #         else:
    #             output.failure(f"Expected a redirect but got {output.status_code}")
