podtato = {
    "parts": [
        {
            "name": "entry",
            "image": "0.2.7",
            "service": {
                "port": 9000,
                "type": "LoadBalancer"
            }
        },
        {
            "name": "hat",
            "image": "0.2.7",
            "service": {
                "port": 9001,
                "type": "ClusterIP"
            }
        },
        {
            "name": "left-leg",
            "image": "0.2.7",
            "service": {
                "port": 9002,
                "type": "ClusterIP"
            }
        },
        {
            "name": "left-arm",
            "image": "0.2.7",
            "service": {
                "port": 9003,
                "type": "ClusterIP"
            }
        },
        {
            "name": "right-leg",
            "image": "0.2.7",
            "service": {
                "port": 9004,
                "type": "ClusterIP"
            }
        },
        {
            "name": "right-arm",
            "image": "0.2.7",
            "service": {
                "port": 9005,
                "type": "ClusterIP"
            }
        }
    ]
}
