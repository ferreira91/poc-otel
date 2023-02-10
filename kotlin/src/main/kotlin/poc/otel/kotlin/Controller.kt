package poc.otel.kotlin

import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@RestController
class Controller {

    @GetMapping("/otel")
    fun otel(): String {
        val observability = listOf("logs", "metrics", "tracing")
        return observability.random()
    }
}